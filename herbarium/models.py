from django.db import models

class Family(models.Model):
    class Meta:
        verbose_name_plural = "families"

    name = models.CharField(max_length=200)
    latin = models.CharField(max_length=200)

class Plant(models.Model):
    name = models.CharField(max_length=200)
    latin = models.CharField(max_length=200)
    family = models.ForeignKey(Family, on_delete=models.CASCADE)

class Variety(models.Model):
    class Meta:
        verbose_name_plural = "varieties"

    name = models.CharField(max_length=200)
    latin = models.CharField(max_length=200)
    plant = models.ForeignKey(Plant, on_delete=models.CASCADE)

    edible = models.BooleanField(default=False)
    needs_trellis = models.BooleanField(default=False)
    needs_bird_netting = models.BooleanField(default=False)
    is_keto = models.BooleanField(default=False)
    native = models.BooleanField(default=False)
    invasive = models.BooleanField(default=False)
    is_cover = models.BooleanField(default=False)
