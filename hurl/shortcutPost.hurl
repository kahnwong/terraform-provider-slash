POST {{SLASH_HOST}}/api/v1/shortcuts
Content-Type: application/json
Authorization: Bearer {{SLASH_ACCESS_TOKEN}}
{
  "name": "foo",
  "link": "https://foo.example.com",
  "title": "Foo"
}
