#!/usr/bin/env python
from WMCore.WMExceptions import WMEXCEPTION, WM_JOB_ERROR_CODES, STAGEOUT_ERRORS
with open('WMJobErrors.txt', 'w') as ostream:
    for key, val in WM_JOB_ERROR_CODES.items():
        ostream.write(("{}: {}\n".format(key, val)))
with open('WMExceptions.txt', 'w') as ostream:
    for key, val in WMEXCEPTION.items():
        ostream.write("{}: {}\n".format(key, val))
with open('StageoutErrors.txt', 'w') as ostream:
    for key, vals in STAGEOUT_ERRORS.items():
        msgs = []
        for item in vals:
            msgs.append('if ({}) {}'.format(item['regex'], item['error-msg']))
        ostream.write("{}: {}\n".format(key, '; '.join(msgs)))
