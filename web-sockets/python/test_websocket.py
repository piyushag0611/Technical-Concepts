import pytest
import asyncio
from server import start_server
import server
import client
from client import connect_to_server

# Server starts and listens on a port.
@pytest.mark.asyncio
async def test_server_starts_and_listens():

    server = await start_server()
    assert server.is_serving()
    server.close()
    await server.wait_closed()
    assert not server.is_serving()

# Client can connect to the server.
@pytest.mark.asyncio
async def test_client_can_connect():

    async def handler(ws):
        await ws.wait_closed()

    server = await start_server(handler)
    async with connect_to_server() as websocket:
        assert websocket.state.name == "OPEN"
    server.close()
    await server.wait_closed()
    assert not server.is_serving()

# On connection, client receives a welcome message from the server.
@pytest.mark.asyncio
async def test_on_connection_server_message_client_receives():

    server = await start_server()
    async with connect_to_server() as websocket:
        message = await asyncio.wait_for(websocket.recv(), timeout=2.0)
        assert message == "Welcome!"
    
    server.close()
    await server.wait_closed()
    assert not server.is_serving()

@pytest.mark.asyncio
async def test_client_send_recieve_echo():

    server = await start_server()
    async with connect_to_server() as websocket:
        message = await client.receive_message(websocket)
        assert message == "Welcome!"

        message = "Hello, how do you do?"
        await client.send_message(websocket, message)
        message_received = await client.receive_message(websocket)
        assert message_received == message

    server.close()
    await server.wait_closed()
    assert not server.is_serving()

# Server sends an independent message, and client receives the message correctly.
@pytest.mark.asyncio
async def test_server_send_message_client_receive():

    async def handler(ws):

        await ws.send("Welcome!")
        await ws.send("PING")
        await ws.wait_closed()

    _server = await start_server(handler)
    async with connect_to_server() as websocket:
        message = await client.receive_message(websocket)
        assert message == "Welcome!"

        server_message = "PING"
        message = await client.receive_message(websocket)
        assert message == server_message
    
    _server.close()
    await _server.wait_closed()
    assert not _server.is_serving()

# Client can close the connection cleanly.
@pytest.mark.asyncio
async def test_close_connection():

    _server = await start_server()
    async with connect_to_server() as websocket:
        message = await client.receive_message(websocket)
        assert message == "Welcome!"

        await websocket.close(code=1000)
        assert websocket.state.name == "CLOSED"
        assert websocket.close_code == 1000
    
    _server.close()
    await _server.wait_closed()
    assert not _server.is_serving()

# Server handles multiple sequential messages correctly.
@pytest.mark.asyncio
async def test_server_sequential_messages():

    _server = await start_server()
    async with connect_to_server() as websocket:
        message = await client.receive_message(websocket)
        assert message == "Welcome!"

        messages = ["hello", "how", "do", "you", "do", "?"]
        recvd_messages = []
        for message in messages:
            await client.send_message(websocket, message)
            recv_message = await client.receive_message(websocket)
            recvd_messages.append(recv_message)
        
        assert recvd_messages == messages

        await websocket.close(code=1000)
        assert websocket.state.name == "CLOSED"
        assert websocket.close_code == 1000
    
    _server.close()
    await _server.wait_closed()
    assert not _server.is_serving()

# Client handles multiple messages correctly
@pytest.mark.asyncio
async def test_client_sequential_messages():

    server_messages = ["I", "am", "doing", "fine"]
    async def handler(ws):

        await ws.send("Welcome!")
        for message in server_messages:
            await ws.send(message)
        await ws.wait_closed()
    
    _server = await start_server(handler)

    async with connect_to_server() as websocket:
        message = await client.receive_message(websocket)
        assert message == "Welcome!"

        server_sent_messages = await client.receive_indp_messages(websocket)
        assert server_sent_messages == server_messages

        await websocket.close(code=1000)
        assert websocket.state.name == "CLOSED"
        assert websocket.close_code == 1000
    
    _server.close()
    await _server.wait_closed()
    assert not _server.is_serving()

# Server handles multiple clients independently
@pytest.mark.asyncio
async def test_multiple_clients():

    _server = await start_server()
    
    async with connect_to_server() as ws1, connect_to_server() as ws2, connect_to_server() as ws3:
        for ws in [ws1, ws2, ws3]:

            message = await client.receive_message(ws)
            assert message == "Welcome!"
        
        messages = [(ws1, "Ping1"), (ws2, "Ping2"), (ws3, "Ping3")]

        tasks = [client.send_message(ws, message) for (ws, message) in messages]
        await asyncio.gather(*tasks)

        for (ws, message) in messages:

            recv_message = await client.receive_message(ws)
            assert recv_message == message

    _server.close()
    await _server.wait_closed()
    assert not _server.is_serving()









