import datetime

from django.db import models

class Box(models.Model):
    class Meta:
        verbose_name_plural = "boxes"
        ordering = ['name']

    name = models.CharField(max_length=200)
    want_qr = models.BooleanField('QR code?', default=False)

    def __str__(self):
        return "{} ({} things)".format(self.name, self.boxcontent_set.count())

class BoxContent(models.Model):
    class Meta:
        ordering = ['when_planted', 'name']

    box = models.ForeignKey(Box, on_delete=models.CASCADE)
    name = models.CharField(max_length=200)
    when_planted = models.PositiveIntegerField(default=2022)
