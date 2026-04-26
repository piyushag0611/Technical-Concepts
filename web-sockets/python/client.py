import websockets
import asyncio

def connect_to_server():

    return websockets.connect("ws://localhost:8765")