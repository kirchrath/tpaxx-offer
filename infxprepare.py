#!/usr/bin/env python
import re
import sys
from datetime import datetime
from dateutil.relativedelta import relativedelta

filepath = sys.argv[1]

pattern = '(?P<v>.{2})(?P<marke>.{4})(?P<outsource>.{3})(?P<outdest>.{3})' \
'(?P<insource>.{3})(?P<indest>.{3})(?P<kennung>.{1})(?P<art>.{4})' \
'(?P<termin>.{10})(?P<dauer>.{10})(?P<code>.{17})(?P<zimmer>.{3})' \
'(?P<verpfl>.{3})(?P<belegung>.{1})(?P<preis>.{5})(?P<airline>.{2})' \
'(?P<artlang>.{40})(?P<hname>.{20})(?P<ort>.{20})(?P<beschr>.{240})' \
'(?P<kate>.{3})(?P<kname>.{40})(?P<kseite>.{4})(?P<restpl>.{6})' \
'(?P<bild>.{120})(?P<typ>.{10})(?P<anreise>.{1})(?P<bmin>.{1})' \
'(?P<bmax>.{1})(?P<vmin>.{1})(?P<vmax>.{1})(?P<currency>.{3})'

create_table_sql = 'CREATE TABLE IF NOT EXISTS offers (id INT NOT NULL PRIMARY KEY,' \
' outsource CHARACTER(3), outdest CHARACTER(3), start DATE, duration INT, ' \
'hotelcode CHARACTER(20), accommodation CHARACTER(3), catering CHARACTER(3), ' \
'carrier CHARACTER(3), operator CHARACTER(3), category INT, ' \
'type CHARACTER(20), bmin INT, bmax INT, vmin INT, vmax INT, ' \
'belegung INT, preis DECIMAL(10,4), currency CHARACTER(3) ); '

clear_table = 'delete from offers;'

insert_row_sql = "insert into offers values({},'{}','{}','{}',{},'{}','{}','{}','{}','{}',{},{},{},{},{},{},{},{},'{}');"

offerId = 0
with open(filepath) as fp:  
  line = fp.readline()
  print("begin transaction;")
  print(create_table_sql)
  print(clear_table)

  while line:
    offerId = offerId + 1
    data = line.strip()
    match = re.search(pattern, data)
    termin = match.group('termin').strip()
    termin_object = datetime.strptime(termin, '%d.%m.%y') + relativedelta(years=+5)

    tmp = insert_row_sql.format(
      offerId,
      match.group('outsource').strip(),
      match.group('outdest').strip(),
      datetime.strftime(termin_object, '%Y-%m-%d'),
      match.group('dauer').strip(),
      match.group('code').strip(),
      match.group('zimmer').strip(),
      match.group('verpfl').strip(),
      match.group('airline').strip(),
      match.group('airline').strip(),
      match.group('kate').strip(),
      match.group('typ').strip(),
      match.group('bmin').strip(),
      match.group('bmax').strip(),
      match.group('vmin').strip(),
      match.group('vmax').strip(),
      match.group('belegung').strip(),
      str(int(match.group('preis').strip())),
      match.group('currency').strip()
    )
    print(tmp)
    line = fp.readline()
  print('commit;')



# offerId text NOT NULL PRIMARY KEY, hotelCode text NOT NULL, duration integer NOT NULL, total text NOT NULL
