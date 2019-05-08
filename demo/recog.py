# It helps in identifying the faces
import cv2, sys, numpy, os, threading
from http.server import BaseHTTPRequestHandler, HTTPServer
import simplejson
import time
import base64
import http.server
import socketserver

size = 4
haar_file = 'haarcascade_frontalface_default.xml'
datasets = 'data'
data_lock = threading.Lock()


class S(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        self._set_headers()
        f = open("index.html", "r")
        self.wfile.write(f.read())

    def do_HEAD(self):
        self._set_headers()

    def do_POST(self):
        self._set_headers()
        print("in post method")
        self.data_string = self.rfile.read(int(self.headers['Content-Length']))

        # Convert back to binary
        png_original = base64.b64decode(self.data_string)

        # Write to a file to show conversion worked
        with open(self.path.repalce('/', '') + "/1" + '.png', 'wb') as f_output:
            f_output.write(png_original)
        # self.send_response(200)
        # self.end_headers()
        # print("data_string", self.data_string)
        # print("got post!!")
        # content_len = int(self.headers.getheader('content-length', 0))
        # post_body = self.rfile.read(content_len)
        # print("post_body", type(post_body), post_body)
        # test_data = simplejson.loads(post_body)
        # print("test data:", type(test_data), test_data)
        #
        # data = simplejson.loads(self.data_string)
        # print("data_string", type(data), data)
        # self.wfile.write(bytes("received", "utf-8"))
        return


def run(server_class=HTTPServer, handler_class=S, port=8000):
    server_address = ('127.0.0.1', port)
    httpd = server_class(server_address, handler_class)
    print('Starting httpd...')
    print(time.asctime(), 'Server started on port', port)
    print('....')
    print('ctrl-c to quit server.')
    close = False
    try:
        httpd.serve_forever()
        if close:
            raise
            # server.server_close()
    except KeyboardInterrupt:
        httpd.server_close()
        print(time.asctime(), "Server Stopped")
    except:
        httpd.server_close()
        print(time.asctime(), "Server Stopped")


# Part 1: Create fisherRecognizer
print('Recognizing Face Please Be in sufficient Lights...')

# Create a list of images and a list of corresponding names
(images, lables, names, id) = ([], [], {}, 0)
for (subdirs, dirs, files) in os.walk(datasets):
    for subdir in dirs:
        names[id] = subdir
        subjectpath = os.path.join(datasets, subdir)
        for filename in os.listdir(subjectpath):
            path = subjectpath + '/' + filename
            lable = id
            images.append(cv2.imread(path, 0))
            lables.append(int(lable))
            print(type(images), type(lable))
        id += 1
(width, height) = (130, 100)


def update_data():
    global id
    global images
    global lables
    global names
    global model
    while True:
        for (subdirs, dirs, files) in os.walk(datasets):
            for subdir in dirs:
                if subdir not in names.values():
                    print(subdir, names)
                    with data_lock:
                        names[id] = subdir
                        subjectpath = os.path.join(datasets, subdir)
                        for filename in os.listdir(subjectpath):
                            path = subjectpath + '/' + filename
                            lable = id
                            image = numpy.array([cv2.imread(path, 0)])
                            print(image.shape)
                            images = numpy.append(images, image)
                            lables = numpy.append(lables, lable)
                            model.update(image, numpy.array(lable))
                            print(type(images), type(lable))
                        id += 1


# def http_server():
#     port = 8000
#
#     handler = http.server.SimpleHTTPRequestHandler
#
#     with socketserver.TCPServer(("", port), handler) as httpd:
#         print("serving at port", port)
#         httpd.serve_forever()
#     #
#     # httpd = SocketServer.TCPServer(("", port), handler)
#     #
#     # print("serving at port", port)
#     # httpd.serve_forever()


# Create a Numpy array from the two lists above
# (images, lables) = [numpy.array(lis) for lis in [images, lables]]
# print(type(images), type(lable))

# OpenCV trains a model from the images
# NOTE FOR OpenCV2: remove '.face'
(images, lables) = [numpy.array(lis) for lis in [images, lables]]
print(type(images), lables, type(lables))
model = cv2.face.LBPHFaceRecognizer_create()
model.train(images, lables)

# run update thread in parallel
update_data = threading.Thread(target=update_data)
update_data.daemon = True
update_data.start()

# run http server in parallel
server_run = threading.Thread(target=run)
server_run.daemon = True
server_run.start()

# Part 2: Use fisherRecognizer on camera stream
face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_default.xml")
webcam = cv2.VideoCapture(0)
while True:
    (_, im) = webcam.read()
    gray = cv2.cvtColor(im, cv2.COLOR_BGR2GRAY)
    faces = face_cascade.detectMultiScale(gray, 1.3, 5)
    for (x, y, w, h) in faces:
        cv2.rectangle(im, (x, y), (x + w, y + h), (255, 0, 0), 2)
        face = gray[y:y + h, x:x + w]
        face_resize = cv2.resize(face, (width, height))
        with data_lock:
            # Try to recognize the face
            prediction = model.predict(face_resize)
            match_name = names[prediction[0]]
        prediction_score = prediction[1]
        cv2.rectangle(im, (x, y), (x + w, y + h), (0, 255, 0), 3)

        if prediction_score < 80:
            cv2.putText(im, '% s - %.0f' % (match_name, prediction_score), (x - 10, y - 10), cv2.FONT_HERSHEY_PLAIN, 1,
                        (0, 255, 0))
        else:
            cv2.putText(im, 'not recognized', (x - 10, y - 10), cv2.FONT_HERSHEY_PLAIN, 1, (0, 255, 0))

    cv2.imshow('OpenCV', im)

    key = cv2.waitKey(10)
    if key == 27:
        break

