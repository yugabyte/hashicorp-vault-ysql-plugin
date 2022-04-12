#   Issues Faced
-   This file contains the list of issues while developing the plugin and what methods were used to resolve them.

-   The list of issues were:
    -   DeleteUser

##  DeleteUser
-   Before Droping the User we are required to remove all the access from him.
-   To do so we can either revoke them or transfer them to some other user.
-   The details are available [here](https://stackoverflow.com/questions/3023583/how-to-quickly-drop-a-user-with-existing-privileges).
-   A user with the access privillage will give error on the Drop user call.
-   https://www.postgresql.org/docs/current/role-removal.html
-   The same error as of [this one](https://www.postgresql.org/message-id/83894A1821034948BA27FE4DAA47427928F7C29922%40apde03.APD.Satcom.Local) is faced when the delete role commands similar to the postgres is implemented for yugabyteDB plugin.
-    When you remove a role referenced in any database, PostgreSQL will raise an error. In this case, you have to take two steps:
    -   First, either remove the database objects owned by the role using the DROP OWNED statement or reassign the ownership of the database objects to another role REASSIGN OWNED.
    -   Second, revoke any permissions granted to the role.