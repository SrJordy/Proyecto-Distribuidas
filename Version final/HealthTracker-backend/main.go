package main

import (
	"GolandProyectos/handlers"
	"GolandProyectos/models"
	"GolandProyectos/repository"
	"GolandProyectos/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	app := fiber.New() // Inicia la aplicación Fiber

	// Configura el middleware CORS para permitir solicitudes de cualquier origen
	app.Use(cors.New())

	config := viper.New() // Inicia Viper para la configuración

	// Configura Viper para leer variables de entorno
	config.AutomaticEnv()
	config.SetDefault("APP_PORT", "3000")
	config.SetDefault("APP_ENV", "development")
	config.SetConfigName("config") // Nombre del archivo de configuración sin la extensión
	config.SetConfigType("env")    // Extensión del archivo de configuración
	config.AddConfigPath(".")      // Ubicación del archivo de configuración
	config.AddConfigPath("/etc/secrets/")

	// Intenta leer el archivo de configuración
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Advertencia: No se pudo leer el archivo de configuración. %v", err)
	}

	// Establece la cadena DSN para la conexión a la base de datos
	dsn := "host=ep-lingering-snowflake-a5j9m53w.us-east-2.aws.neon.tech user=AleVl password=VamLyM2btnd4 dbname=healthTracker port=5432 sslmode=verify-full"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	//NUEVOS CAMBIOS
	// Automigración para el modelo User
	if err := db.AutoMigrate(&models.User{},
		&models.Cuidador{},
		&models.Paciente{},
		&models.PacienteCuidador{},
		&models.Agenda{},
		&models.Medicamento{},
		&models.HorarioMedicamento{},
		&models.HorarioMedicine{},
		&models.Medicine{},
		&models.DeviceToken{}); err != nil {
		log.Fatalf("Error en la automigración: %v", err)
	}

	// Crear instancia del repositorio y handler
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Configurar rutas de usuarios
	routers.SetupUserRoutes(app, userHandler)

	// Crea el repositorio y el handler para PacienteCuidador
	pcRepo := repository.NewPacienteCuidadorRepository(db)   // Asegúrate de haber creado esta función en tu paquete repository
	pcHandler := handlers.NewPacienteCuidadorHandler(pcRepo) // Y esta en tu paquete handlers

	// Configura las rutas de PacienteCuidador
	routers.SetupPacienteCuidadorRoutes(app, pcHandler)

	// Instancia del repositorio y handler para Agenda
	agendaRepo := repository.NewAgendaRepository(db)       // Asegúrate de implementar esto
	agendaHandler := handlers.NewAgendaHandler(agendaRepo) // Y esto también
	routers.SetupAgendaRoutes(app, agendaHandler)          // Asegúrate de implementar SetupAgendaRoutes

	// Instancia del repositorio y handler para Medicamentos
	medicamentoRepo := repository.NewMedicineRepository(db)            // Asegúrate de implementar esto en tu paquete repository
	medicamentoHandler := handlers.NewMedicineHandler(medicamentoRepo) // Y esto en tu paquete handlers
	routers.SetupMedicamentoRoutes(app, medicamentoHandler)            // Incluye esta línea para configurar las rutas de medicamentos

	horarioMedicamentosRepo := repository.NewHorarioMedicamentosRepository(db)                    // Asegúrate de tener esta función implementada
	horarioMedicamentosHandler := handlers.NewHorarioMedicamentosHandler(horarioMedicamentosRepo) // Y esta también
	routers.SetupHorarioMedicamentosRoutes(app, horarioMedicamentosHandler)                       // No olvides implementar esta función

	// Crea la instancia del repositorio de tokens de dispositivos.
	deviceTokenRepo := repository.NewDeviceTokenRepository(db)
	// Crea la instancia del handler de notificaciones.
	notificationHandler := handlers.NewNotificationHandler(deviceTokenRepo)
	// Configura las rutas de notificaciones.
	routers.SetupNotificationRoutes(app, notificationHandler)

	// Define una ruta de bienvenida
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Hola, Mundo!")
	})

	// Inicia el servidor en el puerto configurado
	port := config.GetString("APP_PORT")
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
