# Generated by Django 4.1.7 on 2023-03-02 05:15

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('herbarium', '0010_plantingwindow'),
    ]

    operations = [
        migrations.RenameField(
            model_name='plantingwindow',
            old_name='first_day',
            new_name='end',
        ),
        migrations.RenameField(
            model_name='plantingwindow',
            old_name='last_day',
            new_name='start',
        ),
        migrations.AlterField(
            model_name='plant',
            name='tags',
            field=models.ManyToManyField(blank=True, to='herbarium.tag'),
        ),
    ]
