# Collection of regex patterns

## Match URL and path

```javascript
const pattern =
  /[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)|\/[-a-zA-Z0-9@:%_\+.~#?&//=]*/gi;
```
