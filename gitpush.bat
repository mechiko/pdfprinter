git add . 
IF "%~1" NEQ "" goto Message
git commit -m "default message"
goto End
:Message
git commit -m "%~1"
:End
git push -u origin
@rem git push origin
