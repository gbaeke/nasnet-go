import keras
from keras.applications.nasnet import NASNetMobile
from keras.preprocessing import image
from keras.applications.xception import preprocess_input, decode_predictions
import numpy as np
import tensorflow as tf
from keras import backend as K

sess = tf.Session()
K.set_session(sess)

model = NASNetMobile(weights="imagenet")
img = image.load_img('cat.jpg', target_size=(224,224))
img_arr = np.expand_dims(image.img_to_array(img), axis=0)
x = preprocess_input(img_arr)

preds = model.predict(x)
print('Prediction:', decode_predictions(preds, top=5)[0])

# save the model for use with TensorFlow
builder = tf.saved_model.builder.SavedModelBuilder("nasnet")
# Tag the model, required for Go
builder.add_meta_graph_and_variables(sess, ["atag"])
builder.save()
sess.close()