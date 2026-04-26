import websockets
import asyncio

WELCOME_MESSAGE = "Welcome!"

async def start_server(handler=None):

    if handler is None:
        async def handler(ws):
            await ws.send(WELCOME_MESSAGE)
            async for message in ws:
                await ws.send(message)
        

    server = await websockets.serve(handler, "localhost", 8765)
    return server

async def send_message(ws, message):

    await ws.send(message)

async def receive_message(ws):

    message = await ws.recv()
    await send_message(ws, message)

async def main():

   server = await start_server()
   await server.serve_forever()


if __name__ == "__main__":
    
    asyncio.run(main())