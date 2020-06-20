#!/bin/bash

encodings="utf-8 utf-16le utf-16be utf-32le utf-32be windows-1252 tcvn-5712 cp437 ebcdicatde"


for enc in $encodings; do
  echo "Processing $enc"
  echo -e "This content is encoded as ${enc}\n\n" | iconv -f "UTF-8" -t "${enc}" > $enc.txt
  iconv -f "UTF-8" -t "${enc}//TRANSLIT" ground-truth.txt >> $enc.txt
done
