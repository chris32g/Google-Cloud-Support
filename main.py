# [START gae_python37_app]
from flask import Flask, render_template
from jinja2 import Template
from google.cloud import storage
app = Flask(__name__)


def list_blobs(bucket_name, bucket_folder):
    client = storage.Client()
    #BUCKET_NAME = 'mymadpbucket'
    bucket = client.get_bucket(bucket_name)
    blobs = bucket.list_blobs()

    list = []
    for blob in blobs:
        try:
            num = blob.name.count('/')
            string = blob.name.split('/')[num]
            folder = blob.name.split('/')[0]
            if string != "" and folder == bucket_folder:
                list.append(string)
        except:
            print("An exception occurred")
    return list


@app.route('/')
def home():
    return render_template('main_template.html', my_string="Seleccione la categoria que desea ver!", my_list=['skydive','snowboard'])


@app.route('/skydive')
def skydive():
    return render_template('storage_template.html', my_string="Skydiving Videos!", my_list=list_blobs('mymadpbucket', 'Skydiving'), my_theme='Skydiving')


@app.route('/snowboard')
def snowboard():
    return render_template('storage_template.html', my_string="Snowboarding Videos!", my_list=list_blobs('mymadpbucket', 'Snowboard'), my_theme='Snowboard')


if __name__ == '__main__':
    # This is used when running locally only. When deploying to Google App
    # Engine, a webserver process such as Gunicorn will serve the app. This
    # can be configured by adding an `entrypoint` to app.yaml.
    app.run(host='127.0.0.1', port=8080, debug=True)
# [END gae_python37_app]
