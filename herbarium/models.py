from django.db import models

class Family(models.Model):
    class Meta:
        verbose_name_plural = "families"
        ordering = ['name']

    name = models.CharField(max_length=200)

    def __str__(self):
        return self.name

class Plant(models.Model):
    family = models.ForeignKey(Family, on_delete=models.CASCADE)

    edible = models.BooleanField('Edible?', null=True)
    needs_trellis = models.BooleanField('Needs trellis?', null=True)
    needs_bird_netting = models.BooleanField('Needs bird netting?', null=True)
    is_keto = models.BooleanField('Is keto?', null=True)
    native = models.BooleanField('Native plant?', null=True)
    invasive = models.BooleanField('Is invasive?', null=True)
    is_cover = models.BooleanField('Cover crop?', null=True)
    grow_from_seed = models.BooleanField('Good to grow from seeds?', default=True, null=True)

    class Type(models.TextChoices):
        VEGETABLE = 'V', 'Vegetable'
        FRUIT = 'F', 'Fruit'
        HERB = 'H', 'Herb'
        FLOWER = 'L', 'Flower'
        GREEN = 'G', 'Green'
    
    class Lifespan(models.TextChoices):
        ANNUAL = 'A', 'Annual'
        BIENNIAL = 'B', 'Biennial'
        PERENNIAL = 'P', 'Perennial'
        UNKNOWN = 'U', 'Unknown'

    type = models.CharField(max_length=1, choices=Type.choices)
    lifespan = models.CharField(max_length=1, choices=Lifespan.choices, default=Lifespan.UNKNOWN)
    bad_for_cats = models.BooleanField('Bad for cats?', null=True)
    deer_resistant = models.BooleanField('Deer resistant?', null=True)

    def __str__(self):
        return "{} ({}, {} varieties)".format(self.primary_name(), self.family, self.variety_set.count())

    def primary_name(self):
        return self.plantname_set.all()[0].name

class PlantName(models.Model):
    class Meta:
        verbose_name_plural = "plant names"

    name = models.CharField(max_length=200)
    plant = models.ForeignKey(Plant, on_delete=models.CASCADE)

class Variety(models.Model):
    class Meta:
        verbose_name_plural = "varieties"
        ordering = ['name']

    name = models.CharField(max_length=200)
    plant = models.ForeignKey(Plant, on_delete=models.CASCADE)
    heat_sensitive = models.BooleanField('Cover during heatwaves?', null=True)

    def __str__(self):
        return "{} {} ({})".format(self.name, self.plant.primary_name(), self.plant.family.name)