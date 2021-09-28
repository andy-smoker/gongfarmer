Для начала нужно создать файл config.toml который имеет вид:

```
log = "/log/file/path.log"
work_directory = "/work/ditictory/"
service_dirs = ["1","2","3"] // list of subdirs into workdir for cleaning
ignore_list = ["diriktory", "or", "file", "name"]
login = "atlssian@user.login"
jira_token = "atlassian-token"
hour_of_cleaning = 0  // removal hour 
project = "subdomain-atlassian"
period = 1  // one run of value days
```