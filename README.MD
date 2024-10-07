# Forum

## Description
Crust & Crumbles is a web-based baking forum designed to foster community interaction among baking enthusiasts. It allows users to share and categorize recipes, comment on posts, and like or dislike both posts and comments. The forum features user authentication, session management for secure access, and post filtering to easily navigate through baking tips, recipes, and discussions. Powered by SQLite for data storage, the platform is containerized using Docker for seamless deployment. This project offers a comprehensive exploration of web development, database management, and containerization, tailored to the baking community.

# Consists of the following 
1) communication between users.
2) associating categories to posts.
3) liking and disliking posts and comments.
4) filtering posts.
5) SQLite
6) Authentication
7) SQLite

## Steps to running full server 
1) setup and installation 
2) running sqlite database 
3) running docker interminal/docker 
4) run project 


================================================================

# Setup and Installation
```sh
git clone https://learn.reboot01.com/git/zalhaj/forum
cd forum 
```
# runing sqlite database
to run the sqlite database you need to install a C compiler by running the following commands in poweshell (in your pc) running it as administrators :

1) 
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))


——
2) 
choco install mingw -y


—-
3) 
[System.Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\tools\mingw64\bin", [System.EnvironmentVariableTarget]::Machine)

—-
4) 
go env | findstr /C:"CGO_ENABLED"


# To run the project
```go
go run .
```

# Docker 
To run docker file do the following: 
```
docker image build -f Dockerfile -t dockerize .

docker container run -p 8989:8989 --name forum dockerize
```

# tip
to improve viewing settings of the data base install the vs extension :
- sqlite editor 

================================================================

## AUTHORS  
- Zahra Alhaj (zahraalhaj)
- Dana Alkooheji (dalkooheji)
- Ghadeer Husain (ghadeerl)
- Salman Alkooheji (salkooheji)