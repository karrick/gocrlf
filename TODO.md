# plain

## CRLFfromLF

* LF will change to CRLF
* CRLF will not change
* CR will not change

## LFfromCRLF

* CRLF will change to LF
* CR will not change
* LF will not change

# more practical

## dos2unix

CRLF: LF
CR: ?
LF: LF

## unix2dos

CRLF: CRLF
CR: ?
LF: CRLF

## unix2mac

CRLF: ?
CR: CR
LF: CR

## mac2dos

CRLF: CRLF
CR: CRLF
LF: ?

## mac2unix

CRLF: ?
CR: LF
LF: LF

## dos2mac

CRLF: CR
CR: CR
LF: ?

# even more practical

## 2dos (possibly expanding)

CRLF: CRLF
CR: CRLF
LF: CRLF

## 2mac (possibly reducing)

CRLF: CR
CR: CR
LF: CR

## 2unix (possibly reducing)

CRLF: LF
CR: LF
LF: LF
