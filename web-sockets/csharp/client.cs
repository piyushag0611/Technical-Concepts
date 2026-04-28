using System.Net.WebSockets;
using System.Text;

namespace WebSocketDemo;

public class Client
{
    public async Task<ClientWebSocket> ConnectAsync(string url = "ws://localhost:8765/")
    {
        var ws = new ClientWebSocket();
        await ws.ConnectAsync(new Uri(url), CancellationToken.None);
        return ws;
    }

    public async Task<string> SendMessage(ClientWebSocket conn, string message)
    {
        var bytes = Encoding.UTF8.GetBytes(message);
        await conn.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var echoMessage = Encoding.UTF8.GetString(buffer, 0, result.Count);
        return echoMessage;
    } 

    public async Task<WebSocketState> CloseConnection(ClientWebSocket conn)
    {
        try {
        await conn.CloseAsync(WebSocketCloseStatus.NormalClosure, "", CancellationToken.None);
        return conn.State;
        }
        catch
        {
            Console.WriteLine("Error while closing connection");
            return conn.State;
        }
    }

    public async Task<List<string>> ReceiveServerMessages(ClientWebSocket conn)
    {
        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(500));
        List<string> receivedMessages = [];
        try
        {
            while(true)
            {
                var buffer = new byte[1024];
                var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), cts.Token);
                var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
                receivedMessages.Add(message);
            }
        }
        catch (OperationCanceledException)
        {
            return receivedMessages;
        }
    }
}