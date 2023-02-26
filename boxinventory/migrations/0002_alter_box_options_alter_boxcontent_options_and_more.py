# Generated by Django 4.1.7 on 2023-02-26 19:28

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('boxinventory', '0001_initial'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='box',
            options={'verbose_name_plural': 'boxes'},
        ),
        migrations.AlterModelOptions(
            name='boxcontent',
            options={'ordering': ['when_planted', 'name']},
        ),
        migrations.AddField(
            model_name='box',
            name='want_qr',
            field=models.BooleanField(default=False),
        ),
        migrations.AlterField(
            model_name='boxcontent',
            name='when_planted',
            field=models.PositiveIntegerField(default=2022),
        ),
    ]
