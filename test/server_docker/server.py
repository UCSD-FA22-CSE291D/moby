from http.server import HTTPServer, BaseHTTPRequestHandler
import time

class SlowLorisServer(BaseHTTPRequestHandler):
    def do_GET(self):
        # send 200 response
        self.send_response(200)

        # send response headers
        self.end_headers()

        # Some randomly large memory
        # Single dump size is 613M
        # 1 pre-dump size is 612M + 684K
        garbage = [0] * (256*1024*300)

        # send the body of the response
        self.wfile.write(bytes("Start sending packets...\n", "utf-8"))
        self.wfile.write(bytes("Format: [id] time\n", "utf-8"))
        for i in range(15):
            garbage[i] += 1
            # self.wfile.write(bytes("[%s] %s\n" % (i, datetime.datetime.now().strftime('%X')), "utf-8"))
            self.wfile.write(bytes("[%s] %s\n" % (i, time.time_ns()), "utf-8"))
            self.wfile.flush()
            time.sleep(1)

httpd = HTTPServer(('0.0.0.0', 8000), SlowLorisServer)
httpd.serve_forever()
