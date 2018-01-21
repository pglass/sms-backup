Overview
--------

This is some code to analyze an backup of SMS/MMS messages taken by the SMS
Backup and Restore app for Android.

It can analyze a single conversation and output plots for a few different things:

- Time series of messages sent per day
- Time series of messages sent per week
- Histogram of incoming message lengths
- Histogram of outgoing message lengths
- Scatterplot over time, plotting the hour of day when messages are sent

Other things / TODO
-------------------

- How many emojis or pictures used?
- Most common words and phrases
- Number of exclamations? Number of questions?

Setup
-----

This is written in Go.

I use [glide](https://github.com/Masterminds/glide) to manager go dependencies.
Glide fetches your dependencies into the vendor directory (so you do not have
to commit the vendor directory)

```bash
$ glide install
```

### Build the code

Use the Makefile to build the code:

```bash
$ make main
```

This outputs an executable called `main`:

```bash
$ ./main -h
Usage of ./main:
  -f string
    	The XML file containing your SMS backups
  -n string
    	My phone number. Used to determine if MMS messages are incoming
  -o string
    	The output image (default "out.png")
  -t string
    	One of: messagesPerDay, messagesPerWeek, incomingMessageLengths, outgoingMessageLengths, messagesTimeOfDay
```

How To
------

You will need an XML backup file _for a single conversation_ from the [SMS
Backup & Restore](https://play.google.com/store/apps/details?id=com.riteshsahu.SMSBackupRestore)
Android app. This tool is designed to analyze backups of single conversion,
from a single contact. It may "run" with multi-conversation backups but does
not perform any multi-conversation analysis.

With the current version (as of 1.20.2018), the backup files are named like
`sms-20180119165444.xml`. I uploaded an XML backup to Google Drive from my
phone and then downloaded the file to a laptop. The app supports toggles for
optionally backing up images and emojis - this code supports those backups with
images and emojis.

See the Makefile for more.
