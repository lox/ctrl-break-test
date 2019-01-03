FROM mcr.microsoft.com/windows/servercore:ltsc2016
COPY ./ctrl-break-test.exe /
CMD ["/ctrl-break-test.exe"]