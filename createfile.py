#!/usr/bin/python

import sys,os,math

value=""

mydict={}


def update_dict(key,value):
                mydict.update({key : value })



for line in reversed(list(open("log.txt"))):
	value=line.rstrip()
	break	

ret=value.find('Server is shutting down...All client connections are closed')
if(ret==-1):
	for line in list(open("log.txt")):
		a=line.find('Set ---> Key:')
		if(a!=-1):
			newline=line.split('\n')
			spaceline=newline[0].split(' ')
			keyindex=spaceline.index("Key:")
			valueindex=spaceline.index("Value:")
			#print(spaceline[keyindex+1]+"   "+str(' '.join(spaceline[valueindex+1:])))
			update_dict(spaceline[keyindex+1],str(' '.join(spaceline[valueindex+1:])))	
		a=line.find('Delete ---> Key:')
		if(a!=-1):
			newline=line.split('\n')
			spaceline=newline[0].split(' ')
			keyindex=spaceline.index("Key:")
			#print(spaceline[keyindex+1])
			#del mydict[spaceline[keyindex+1]]	
			mydict.pop(spaceline[keyindex+1], None)
		a=line.find('Rename ---> OldKey:')
		if(a!=-1):
			newline=line.split('\n')
			spaceline=newline[0].split(' ')
			oldkeyindex=spaceline.index("OldKey:")
			newkeyindex=spaceline.index("NewKey:")
			val=mydict[spaceline[oldkeyindex+1]]
			del mydict[spaceline[oldkeyindex+1]]
			#print(str(spaceline[newkeyindex+1]) + str(val) )
			mydict.update({str(spaceline[newkeyindex+1]) : str(val) })	

	theKeys = mydict.keys()
	theKeys.sort()
	dfile=open("input.txt",'w')
	for eachkey in theKeys:		
		dfile.write(str(eachkey)+" "+str(mydict[eachkey])+"\n")
	dfile.close()














