import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service(#{ title: "Calendar API" })
namespace Calendar;

model Health { status: "ok"; }
model BookingInput {
  @minLength(1) title: string;
  startsAt: utcDateTime;
  endsAt: utcDateTime;
}
model Booking extends BookingInput { id: string; }
@error
model BadRequest {
  @statusCode statusCode: 400;
  @body body: ApiError;
}
@error
model NotFound {
  @statusCode statusCode: 404;
  @body body: ApiError;
}
model ApiError { code: string; message: string; }

@route("/health") @get @operationId("getHealth")
op getHealth(): Health;

@route("/api/bookings") @post @operationId("createBooking")
op createBooking(@body booking: BookingInput): {
  @statusCode statusCode: 201;
  @body body: Booking;
} | BadRequest;

@route("/api/bookings/{id}") @get @operationId("getBooking")
op getBooking(@path id: string): Booking | BadRequest | NotFound;
