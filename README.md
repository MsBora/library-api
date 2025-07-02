# library-api
Using Docker so as not to install and configure PostgreSQL directly on the system. The database runs in an isolated, self-contained environment called a container. This is great for a few reasons:
    Consistency: The database runs the same way on your machine, your teammate's machine, or a server.
    
    Isolation: It won't interfere with any other software or databases on your system.
    
    Easy Cleanup: When you're done, you can stop and remove the container, and your system is left perfectly clean.

command:
docker run --name library-db -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres