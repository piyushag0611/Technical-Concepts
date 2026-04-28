import websockets
import asyncio


def connect_to_server():

    return websockets.connect("ws://localhost:8765")

async def send_message(ws, message):

    await ws.send(message)

async def receive_message(ws):

    message = await asyncio.wait_for(ws.recv(), timeout=2.0)
    return message



async def receive_indp_messages(ws, timeout=0.5):
    received = []
    try:
        while True:
            message = await asyncio.wait_for(ws.recv(), timeout=timeout)
            received.append(message)
    except asyncio.TimeoutError:
        return received

async def main():
    async with connect_to_server() as ws:
        welcome = await receive_message(ws)
        print(welcome)

        while True:
            text = input("You: ")
            if text == "exit":
                break
            await send_message(ws, text)
            echo = await receive_message(ws)
            print(f"Server: {echo}")

if __name__ == "__main__":
    asyncio.run(main())
