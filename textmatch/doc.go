// Package textmatch es el motor de comparación de respuestas en cascada del
// ecosistema EduGo (ADR 0035, plan 045). Consolida en un solo lugar —neutral y
// puro, sin infraestructura ni LLM— las utilidades de similitud que hoy están
// duplicadas y divergentes entre learning y worker.
//
// Se organiza en dos niveles ortogonales:
//
//   - Nivel 1 (Comparator/Cascade): ¿este esperado ≈ este candidato? Una cascada
//     de estrategias baratas→caras (Exact → Fuzzy → … → LLM inyectable). Positivo
//     corta; incierto/negativo escala; un error se propaga.
//   - Nivel 2 (SetMatcher): matchea un CONJUNTO de candidatos contra un CONJUNTO
//     de esperados, con una política de completitud (Strict/Lenient) que es
//     decisión de negocio, ortogonal a cómo se compara un par.
//
// El paquete NO decide negocio: qué significa un no-match (ítem ausente → red del
// profesor) lo interpreta el caller. La estrategia LLM NO vive aquí: depende de
// infra/red, así que se define donde vive el provider y se inyecta en la cascada
// (DIP). Todo el código es determinista y libre de estado salvo la construcción
// explícita de Cascade/SetMatcher.
package textmatch
