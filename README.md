Overview
--------

This is some code to analyze an backup of SMS/MMS messages taken by the SMS
Backup and Restore app for Android.

What I want to know
-------------------

- How many messages sent Each day? Each week? Each month?
- How does time of day, day of week affect the number of messages sent?
- How many emojis used?
- How many pictures sent?
- Length of texts? Number of exclamations? Number of questions?
- Most common words and phrases

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
        One of: messagesPerDay, messagesPerWeek
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
