# Generated by Django 4.1.7 on 2023-02-28 07:56

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('herbarium', '0005_primary_name'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='plant',
            options={'ordering': ['name']},
        ),
    ]
