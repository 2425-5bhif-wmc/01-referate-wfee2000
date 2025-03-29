package at.htl.grpc;

import at.htl.grpc.ChatOuterClass.ClaimNameRequest;
import at.htl.grpc.ChatOuterClass.IncomingMessage;
import at.htl.grpc.ChatOuterClass.OutgoingMessage;
import io.grpc.Metadata;
import io.grpc.Metadata.Key;
import io.quarkus.grpc.GrpcClient;
import io.quarkus.grpc.GrpcClientUtils;
import io.quarkus.runtime.QuarkusApplication;
import io.quarkus.runtime.annotations.QuarkusMain;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.subscription.MultiEmitter;
import java.io.IOException;
import java.util.Scanner;
import org.jline.terminal.Terminal;
import org.jline.terminal.TerminalBuilder;

@QuarkusMain
public class ChatMain implements QuarkusApplication {

    // tag::message[]
    private static final StringBuilder message = new StringBuilder();

    // end::message[]

    // tag::injectStub[]
    // Inject gRPC stub
    @GrpcClient
    Chat chat;

    // end::injectStub[]

    // tag::printMessages[]
    /**
     * Print incoming messages without overwriting the current input line.
     *
     * @param incomingStream stream with incoming messages
     */
    public static void printMessages(Multi<IncomingMessage> incomingStream) {
        incomingStream
            .subscribe()
            .with(incomingMessage ->
                System.out.printf(
                    "\033[2K\r%s: %s\n\rWrite message: %s",
                    incomingMessage.getName(),
                    incomingMessage.getResponse(),
                    message
                )
            );
    }

    // end::printMessages[]

    // tag::sendMessages[]
    /**
     * Read input from CLI and send it to service.
     *
     * @param emitter the emitter that sends messages back to the service
     */
    public static void sendMessages(
        MultiEmitter<? super OutgoingMessage> emitter
    ) {
        try (Terminal terminal = TerminalBuilder.terminal()) {
            // strip flags from terminal
            terminal.enterRawMode(); // <1>

            while (true) {
                try {
                    // fill the message with user input
                    readMessage(terminal); // <2>
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }

                // send the message to service
                emitter.emit(
                    // <3>
                    OutgoingMessage.newBuilder()
                        .setMessage(message.toString())
                        .build()
                );

                // print finished message line
                System.out.printf("\033[2K\rYou wrote: %s\n\r", message); // <4>

                // empty the message
                message.setLength(0); // <5>
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    // end::sendMessages[]

    // tag::readMessage[]
    /**
     * Read string from CLI without control characters.
     *
     * @throws IOException if the connection to the terminal cannot be established
     */
    private static void readMessage(Terminal terminal) throws IOException {
        System.out.print("Write message: ");

        int controlCharCounter = 0;

        while (true) {
            // read one character from the console
            char character = (char) terminal.reader().read();

            if (
                (controlCharCounter == 2 && character == '[') ||
                controlCharCounter == 1
            ) {
                controlCharCounter--;
                continue;
            }

            // stop reading and return
            if (character == '\n' || character == '\r') {
                return;
            }

            // delete single character
            if (character == 127 && !message.isEmpty()) {
                System.out.print("\b \b");
                message.deleteCharAt(message.length() - 1);
            }

            controlCharCounter = 0;

            if (character == 27) {
                controlCharCounter = 2;
            }

            // do not display/send control characters
            if (Character.isISOControl(character)) {
                continue;
            }

            // print character to console and append it to message
            System.out.print(character);

            message.append(character);
        }
    }

    // end::readMessage[]

    // tag::getToken[]
    @Override
    public int run(String... args) {
        // claim name from server
        String token = chat
            .claimName(
                ClaimNameRequest.newBuilder()
                    .setName(String.join(" ", args))
                    .build()
            )
            .await()
            .indefinitely()
            .getToken();
        // end::getToken[]

        // tag::makeHeaders[]
        // append token in header
        Metadata headers = new Metadata();
        headers.put(
            Key.of("authorization", Metadata.ASCII_STRING_MARSHALLER), // <1>
            String.format("Bearer %s", token) // <2>
        );
        // end::makeHeaders[]

        // tag::bindHeaders[]
        // attach headers to stub
        Chat authorizedStub = GrpcClientUtils.attachHeaders(chat, headers);
        // end::bindHeaders[]

        // tag::createMulti[]
        // create output stream
        Multi<OutgoingMessage> outgoingStream = Multi.createFrom()
            .<OutgoingMessage>emitter(ChatMain::sendMessages);
        // end::createMulti[]

        // tag::connect[]
        // connect to service and start printing incoming messages
        printMessages(authorizedStub.connect(outgoingStream));
        // end::connect[]

        return 0;
    }
}
