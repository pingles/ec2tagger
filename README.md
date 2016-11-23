# EC2 Tagger

Bulk add tags to your EC2 instances.

## Usage

```
$ ec2tagger --name=ecs ECS True
```

Would add a tag with the key `ECS` and value `True` to all instances that have a `Name` that contains `ecs`.