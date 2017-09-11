
import csv

#todo - pick these up from command line

_inputfile = "redmine_issues.csv"
_outputfile = "scrumdo.csv"
_projectkey = "NE"
_startindex = 101
_maxrecords  = 10000


_scrumdokeys = ["Card ID","Summary","Detail","Points","Estimated Minutes","Business Value","Assignee",
        "Time Criticality","Risk Reduction / Opportunity Enablement","Cell","Rank","Tags","Collections","Labels",
        "Tasks","Due Date","Created","Modified"]

print("Redmine to  Scrumdo")
print("%s --> %s" % ( _inputfile, _outputfile))
data = csv.DictReader(open(_inputfile))

items = []
index = _startindex
records = 0
for row in data:
    item= {}
    for key in _scrumdokeys: item[key] = ""
    item["Card ID"] = _projectkey + "-" +str(index)
    item["Summary"] = row['#'] + " : " + row["Tracker"] + " : " + row["Subject"]
    item["Detail"] = row["Description"]
    item["Created"]  = row["Created"]
    item["Assignee"] = row["Assignee"]
    item["Tags"]   = row["Category"][3:] #<-- unique to NETS
    items.append(item)
    index += 1
    records = records + 1
    _maxrecords -= 1
    if _maxrecords < 1 : break

with open('scrumdo.csv', 'w') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=_scrumdokeys, quotechar='"',delimiter=",",quoting=csv.QUOTE_ALL)
        writer.writeheader()
        for item in items:
            writer.writerow(item)

print("%s records processed." % records )