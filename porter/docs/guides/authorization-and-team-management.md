Porter supports setting basic authorization permissions via for other members in a Porter project. At the moment, there are 3 roles that can be assigned in a Porter project:

- **Admin:** read/write access to all resources, ability to delete the project and manage team members.
- **Developer:** read/write access to applications, jobs, environment groups, cluster data, and integrations.
- **Viewer:** read access to applications, jobs, environment groups, and cluster data.

# Adding Collaborators

To add a new collaborator to a Porter project, you must be logged in with an **Admin** role. As an admin, you will see a **Settings** tab in the sidebar. Navigate to **Settings** and input the email of the user you would like to add. This will generate an invitation link for the user, which expires in 24 hours. The user will get an email to join the Porter project, but if the email is not delivered, you can copy the invite link and send it to them directly.

![image](https://user-images.githubusercontent.com/23369263/125147098-b00f3100-e0ff-11eb-8579-cc28c1a0badc.png)

> 🚧
> 
> If the user does not have a Porter account, they will be asked to register. After registering, if they are not automatically added to the project, the user should **click the invite link again**.  

# Changing Collaborator Permissions

To change an invite or collaborator role, you must be logged in with an **Admin** role. As an admin, you will se a **Settings** tab in the sidebar. Navigate to **Settings** and lookup on the table the invite/collaborator that you want to change it's role, then click the icon with three dots on the row. This will open a pop up that will allow you to select the new role for that invite/collaborator.

You will note that the user that created the project will not be displayed on the table, and you cannot change your own permissions.

![image](https://user-images.githubusercontent.com/23369263/125147141-ea78ce00-e0ff-11eb-9e8b-a3f126874d12.png)

![image](https://user-images.githubusercontent.com/23369263/125147157-0aa88d00-e100-11eb-8d78-1cf34397cd26.png)

# Removing Collaborators

To remove an invite or a collaborator, you must be logged in with an **Admin** role. As an admin, you will se a **Settings** tab in the sidebar. Navigate to **Settings** and lookup on the table the invite/collaborator that you want to remove then click the trash icon to remove the user from the project or delete the invite.

![image](https://user-images.githubusercontent.com/23369263/125147206-3d528580-e100-11eb-9a58-51885ab8b298.png)
