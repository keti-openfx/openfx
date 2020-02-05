all: Go Python NodeJS Ruby CPP JAVA
  
Go:
	cd go; $(MAKE)
Python:
	cd python; $(MAKE)
NodeJS:
	cd nodejs; $(MAKE)
Ruby:
	cd ruby; $(MAKE)
CPP:
	cd cpp; $(MAKE) docker-build; $(MAKE) docker-push
JAVA:
	cd java; $(MAKE)
git:
	git add .
	git commit -m "$m"
	git push origin master
