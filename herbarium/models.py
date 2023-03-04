import datetime

from django.db import models

class Family(models.Model):
    class Meta:
        verbose_name_plural = "families"
        ordering = ['name']

    name = models.CharField(max_length=200)

    def __str__(self):
        return self.name

class Tag(models.Model):
    name = models.CharField(max_length=200, unique=True)

    def __str__(self):
        return self.name

class Plant(models.Model):
    class Meta:
        ordering = ['name']

    name = models.CharField(max_length=200, default='')
    family = models.ForeignKey(Family, on_delete=models.CASCADE)
    tags = models.ManyToManyField(Tag, blank=True)

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
        return self.name
    
    def set_windows_from_str(self, value):
        ws = []
        for w in value.splitlines():
            window = PlantingWindow(plant=self, short_str=w)
            ws.append(window)
        for w in ws:
            w.save()
        self.plantingwindow_set.set(ws)
        self.save()

class PlantName(models.Model): # Actually an alias
    class Meta:
        verbose_name_plural = "plant names"

    name = models.CharField(max_length=200)
    plant = models.ForeignKey(Plant, on_delete=models.CASCADE)

    def __str__(self):
        return self.name

class PlantingWindowManager(models.Manager):
    def on_date(self, date):
        date = date.replace(year=1900)
        first = date.replace(month=1, day=1)
        days = (date - first).days+1
        return self.filter(start__lte=days, end__gte=days)

class PlantingWindow(models.Model):
    class Type(models.TextChoices):
        INDOORS = 'I', 'Indoors'
        DIRECT = 'D', 'Direct'
        TRANSPLANT = 'T', 'Transplant'
        DIRECT_COVERED = 'C', 'Direct, covered'
        TRANSPLANT_COVERED = 'E', 'Transplant, covered'

    objects = PlantingWindowManager()

    plant = models.ForeignKey(Plant, on_delete=models.CASCADE)    
    type = models.CharField(max_length=1, choices=Type.choices)
    start = models.SmallIntegerField()
    end = models.SmallIntegerField()

    @property
    def type_label(self):
        return PlantingWindow.Type(self.type).label

    @property
    def short_str(self):
        return f'{self.type} {self.start_short_str} {self.end_short_str}'

    @short_str.setter
    def short_str(self, value):
        self.type, start, end = value.split()
        self.start_short_str = start
        self.end_short_str = end

    @property
    def start_short_str(self):
        d = datetime.datetime.now().replace(month=1, day=1) + datetime.timedelta(days=self.start-1)
        return f'{d.month}-{d.day}'
    
    @start_short_str.setter
    def start_short_str(self, value):
        d = datetime.datetime.strptime(value, '%m-%d')
        self.start = (d - datetime.datetime(year=1900, month=1, day=1)).days+1

    @property
    def end_short_str(self):
        d = datetime.datetime.now().replace(month=1, day=1) + datetime.timedelta(days=self.end-1)
        return f'{d.month}-{d.day}'
    
    @end_short_str.setter
    def end_short_str(self, value):
        d = datetime.datetime.strptime(value, '%m-%d')
        self.end = (d - datetime.datetime(year=1900, month=1, day=1)).days+1

    @property
    def start_str(self):
        d = datetime.datetime.now().replace(month=1, day=1) + datetime.timedelta(days=self.start-1)
        return d.strftime('%B %d')
    
    @property
    def end_str(self):
        d = datetime.datetime.now().replace(month=1, day=1) + datetime.timedelta(days=self.end-1)
        return d.strftime('%B %d')