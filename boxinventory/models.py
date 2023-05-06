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

    @property
    def contents_by_year(self):
        ret = {}
        for content in self.boxcontent_set.all():
            if content.planted.year not in ret:
                ret[content.planted.year] = []
            ret[content.planted.year].append(content)
        for year in ret:
            ret[year].sort(key=lambda x: x.name)
        return ret

class BoxContent(models.Model):
    class Meta:
        ordering = ['planted', 'name']

    box = models.ForeignKey(Box, on_delete=models.CASCADE)
    name = models.CharField(max_length=200)
    planted = models.DateField(default=datetime.date.today)
