# Generated by Django 4.1.7 on 2023-03-02 04:02

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('herbarium', '0009_tag_plant_tags'),
    ]

    operations = [
        migrations.CreateModel(
            name='PlantingWindow',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('type', models.CharField(choices=[('I', 'Indoors'), ('D', 'Direct'), ('T', 'Transplant'), ('C', 'Direct, covered'), ('E', 'Transplant, covered')], max_length=1)),
                ('first_day', models.SmallIntegerField()),
                ('last_day', models.SmallIntegerField()),
                ('plant', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='herbarium.plant')),
            ],
        ),
    ]
