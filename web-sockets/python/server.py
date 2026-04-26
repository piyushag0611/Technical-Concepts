import websockets
import asyncio

async def start_server():

    server = await websockets.serve(lambda ws: None, "localhost", 8765)
    return server

async def main():

   server = await start_server()
   await server.serve_forever()


if __name__ == "__main__":
    
    asyncio.run(main())