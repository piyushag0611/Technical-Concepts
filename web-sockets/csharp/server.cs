using System.Net;
using System.Net.WebSockets;
using System.Text;

namespace WebSocketDemo;



public class Server
{
    public const string WELCOME_MESSAGE = "Hello!";
    private readonly HttpListener httpListener = new();
    public bool IsListening => httpListener.IsListening;

    private Func<WebSocket, Task>? _handler;

    public Task StartAsync(string url = "http://localhost:8765/",
                           Func<WebSocket, Task>? handler = null)
    {
        _handler = handler ?? DefaultHandler;
        httpListener.Prefixes.Add(url);
        httpListener.Start();
        _ = AcceptLoopAsync();
        return Task.CompletedTask;
    }

    private async Task AcceptLoopAsync()
    {
        while (IsListening)
        {
            try
            {
                var context = await httpListener.GetContextAsync();
                var wsContext = await context.AcceptWebSocketAsync(subProtocol: null);
                _ = _handler!(wsContext.WebSocket);
            }
            catch (HttpListenerException)
            {
                break;
            }
        }
    }

    private async Task DefaultHandler(WebSocket  ws)
    {
      
        var bytes = Encoding.UTF8.GetBytes(WELCOME_MESSAGE);
        await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
        try {
            while (true)
            {
                var buffer = new byte[1024];
                var result = await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
                if (result.MessageType == WebSocketMessageType.Close) 
                {
                    await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "", CancellationToken.None);
                    Console.WriteLine("Client closed connection cleanly.");
                    break;
                }
                var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
                if (!string.IsNullOrEmpty(message))
                {
                    bytes = Encoding.UTF8.GetBytes(message);
                    await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
                }
            }
        } catch (WebSocketException)
        {
            Console.WriteLine("Client disconnected abruptly.");
        }
    }

    public void Stop()
    {
        httpListener.Stop();
    }


}