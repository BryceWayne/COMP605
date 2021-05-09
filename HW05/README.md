# COMP605/HW05

## Asignment details

* You must use CUDA to convert an image in RGB color model to grayscale

* An RGB image can be manipulated as an array that contains 3bytes (char) for each pixel, where each byte represents the intensity of each color

* You must use the *luminosity method* to attain the effect 

$$\texttt{NewImage = (0.3 * R) + (0.59 * G) + (0.11 * B)}$$

* You must produce a report, and compare using different block and grid sizes

* First step is to read a .jpg image and convert it to a 1Darray of uchar4

* Then, write a kernel that uses 2D grid of 2D blocks to manipulate the array

* Be careful  with  the  indexing  (local indices  of  thread,  and  global  indices of  input/output arrays)

* Convert grayscale array to a .jpg image and save it to disk

**NOTE:** You can use OpenCV to read and write from/to .jpg

## HowTo

I wrote the solution in Go.
The code is modified from a blog post that had code which did not run out of the box.
The bug was remedied rather easily after reading a Stack Overflow question.
My edits amount to changing the values according to the *luminosity method*.

```bash
$ ./HW05 filename.jpg
```

**Example:** This will convert a color picture of batman, thus `batman.jpg`, into a graysacale picture.
```bash
$ ./HW05 batman.jpg
```

