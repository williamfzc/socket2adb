import socket
import pprint


# conf
HOST = '127.0.0.1'
PORT = 5037
DEFAULT_ENCODING = 'utf-8'

# connect to adb server
client = socket.socket()
host = HOST or '127.0.0.1'
port = PORT or 5037
client.connect((host, port))


# message we sent to server contains 2 parts
# 1. length
# 2. content
def encode_data(data):
    byte_data = data.encode(DEFAULT_ENCODING)
    byte_length = "{0:04X}".format(len(byte_data)).encode(DEFAULT_ENCODING)
    # looks like
    # b'000Chost:devices'
    return byte_length + byte_data


# ok, we sent it to adb server
# here is an example to get all connected devices
# all services were provided here:
# https://android.googlesource.com/platform/system/core/+/jb-dev/adb/SERVICES.TXT
ready_data = encode_data('host:devices')
client.send(ready_data)


# and, message we got also contains 3 parts:
# 1. status (4)
# 2. length (4)
# 3. content (unknown)
def read_all_content(target_socket):
    result = b''
    while True:
        buf = target_socket.recv(1024)
        if not len(buf):
            return result
        result += buf


# get them
status = client.recv(4)
length = client.recv(4)
content = read_all_content(client)

# check your result
final_result = {
    'status': status,
    'length': length,
    'content': content,
}
final_result = {_: v.decode(DEFAULT_ENCODING) for _, v in final_result.items()}
pprint.pprint(final_result)

# close socket after usage
client.close()
