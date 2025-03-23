package at.htl.grpc;

import io.quarkus.grpc.GrpcService;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import java.time.Duration;
import java.util.concurrent.atomic.AtomicReference;

@GrpcService
public class GreeterService implements Greeter {

    // tag::sayHello[]
    @Override
    public Uni<HelloReply> sayHello(HelloRequest request) {
        return Uni.createFrom()
            .item(
                HelloReply.newBuilder() // <1>
                    .setMessage(greet(request.getName())) // <2>
                    .build() // <3>
            );
    }
    // end::sayHello[]

    // tag::streamHello[]
    public Multi<HelloReply> streamHello(Multi<HelloRequest> incomingStream) {
        AtomicReference<String> name = new AtomicReference<>(""); // <1>
        incomingStream // <2>
            .subscribe()
            .with(request -> {
                name.set(request.getName());
            });
        return Multi.createFrom() // <3>
            .ticks()
            .every(Duration.ofSeconds(1))
            .map(ignored ->
                HelloReply.newBuilder().setMessage(greet(name.get())).build()
            );
    }
    // end::streamHello[]

    private String greet(String name) {
        return String.format("Hello %s!", name);
    }
}
