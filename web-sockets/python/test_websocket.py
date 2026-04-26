import pytest
import asyncio
from server import start_server

@pytest.mark.asyncio

async def test_server_starts_and_listens():

    server = await start_server()
    assert server.is_serving()
    server.close()
    await server.wait_closed()
    assert not server.is_serving()