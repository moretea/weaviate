import sys

# load old header
with open("./tools/header-old.txt", 'r') as oldHeaderFile:
    oldHeader = oldHeaderFile.read()

# load file to edit
with open(sys.argv[1], 'r') as fileToEdit:
    newFile = fileToEdit.read()

# replace the old header
print(newFile.replace(oldHeader, ''))