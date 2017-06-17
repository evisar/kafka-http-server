import asyncio
from aiohttp import web
from aiohttp.web import Application, Response
from aiohttp_sse import sse_response


async def hello(request):
    loop = request.app.loop
    resp = await sse_response(request)
    async with resp:
        for i in range(0, 100):
            print('foo')
            await asyncio.sleep(1, loop=loop)
            resp.send('foo {}'.format(i))
    return resp


async def index(request):
    d = """
        <html>
        <head>
            <script type="text/javascript"
                src="http://code.jquery.com/jquery.min.js"></script>
            <script type="text/javascript">
            var evtSource = new EventSource("/hello");
            evtSource.onmessage = function(e) {
             $('#response').html(e.data);
            }

            </script>
        </head>
        <body>
            <h1>Response from server:</h1>
            <div id="response"></div>
        </body>
    </html>
    """
    return Response(text=d, content_type='text/html')


app = web.Application()
app.router.add_route('GET', '/hello', hello)
app.router.add_route('GET', '/index', index)
web.run_app(app, host='127.0.0.1', port=8080)