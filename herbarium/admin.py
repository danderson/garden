from django.contrib import admin

from .models import Family, Plant, Variety, PlantName

class VarietyInline(admin.StackedInline):
    model = Variety
    extra = 1

class PlantNameInline(admin.StackedInline):
    model = PlantName
    extra = 1

class PlantAdmin(admin.ModelAdmin):
    inlines = [PlantNameInline, VarietyInline]
    fieldsets = [
        (None, {
            'fields': ('family',
                       'type',
                       'lifespan'),
        }),
        ('Extra', {
            'classes': ('collapse',),
            'fields': ('edible',
                       'needs_trellis',
                       'needs_bird_netting',
                       'is_keto',
                       'native',
                       'invasive',
                       'is_cover',
                       'grow_from_seed',
                       'bad_for_cats',
                       'deer_resistant'),
        }),
    ]
    save_on_top = True

admin.site.register(Family)
admin.site.register(Plant, PlantAdmin)
admin.site.register(Variety)
