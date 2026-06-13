# crimes-loader

Repositorio con dos flujos Go integrados:

- `data/`: carga, limpieza y generación de análisis del dataset de Chicago Crimes.
- `ml/`: entrenamiento de un Random Forest sobre los registros limpios.

## Requisitos

- Docker
- Docker Compose

## Estructura esperada de datos

El flujo principal usa estos archivos dentro de `data/`:

- `data/crimes.csv`: dataset fuente.
- `data/output/clean_records_sample.json`: salida del loader.
- `data/output/train_results.json`: salida del trainer.

Si `data/crimes.csv` no existe, el loader genera un dataset sintético de prueba.

## Ejecutar todo el flujo

Desde la raíz del repositorio:

```bash
docker compose up loader trainer
```

Esto:

1. Construye y ejecuta el loader.
2. Genera los JSON en `data/output/`.
3. Ejecuta el trainer usando `clean_records_sample.json`.

## Construir imágenes

```bash
docker compose build loader trainer
```

## Ejecutar solo el trainer

Útil si ya existe `data/output/clean_records_sample.json`:

```bash
docker compose run --rm trainer
```

## Ver resultados

Los archivos generados quedan en:

- `data/output/load_stats.json`
- `data/output/risk_analysis.json`
- `data/output/clean_records_sample.json`
- `data/output/train_results.json`

## Notas

- El compose raíz es el punto de entrada principal.
- El directorio `ml/` conserva su `Dockerfile` propio para el trainer, pero no necesita un `docker-compose.yml` aparte.