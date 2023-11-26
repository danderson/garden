import datetime

from django.db import models

class Box(models.Model):
    class Meta:
        verbose_name_plural = "boxes"
        ordering = ['name']

    name = models.CharField(max_length=200)
    want_qr = models.BooleanField('QR code?', default=False)
    qr_applied = models.BooleanField('QR code applied?', default=False)

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
    name = models.CharField(max_length=200, blank=False)
    latin_name = models.CharField(max_length=200, default="", blank=True)
    planted = models.DateField(default=datetime.date.today)
    removed = models.DateField(default=None, null=True, blank=True)

    def _pretty_date(self, date):
        if date.month == 1 and date.day == 1:
            return '~{}'.format(date.year)
        return date.isoformat()

    @property
    def planted_pretty(self):
        return self._pretty_date(self.planted)

    @property
    def removed_pretty(self):
        return self._pretty_date(self.removed)

    @property
    def dated_name(self):
        paren = 'planted {}'.format(self.planted_pretty)
        if self.removed:
            paren = '{}, removed {}'.format(paren, self.removed_pretty)
        if self.latin_name:
            paren = '{}, {}'.format(self.latin_name, paren)
        return '{} ({})'.format(self.name, paren)
