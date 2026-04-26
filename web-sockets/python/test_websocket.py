import pytest
import asyncio
from server import *
from client import *

@pytest.mark.asyncio
async def test_server_starts_and_listens():

    server = await start_server()
    assert server.is_serving()
    server.close()
    await server.wait_closed()
    assert not server.is_serving()

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


