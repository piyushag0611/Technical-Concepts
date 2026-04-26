import websockets
import asyncio

async def start_server(handler=None):

    if handler is None:
        async def handler(ws):
            await ws.wait_closed()

    server = await websockets.serve(handler, "localhost", 8765)
    return server

async def main():

   server = await start_server()
   await server.serve_forever()


if __name__ == "__main__":
    
    asyncio.run(main())