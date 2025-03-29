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
import org.jline.terminal.Terminal;
import org.jline.terminal.TerminalBuilder;

@QuarkusMain
public class ChatMain implements QuarkusApplication {

    private static final StringBuilder message = new StringBuilder();

    // TODO: Inject gRPC stub

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
            terminal.enterRawMode();

            while (true) {
                try {
                    // fill the message with user input
                    readMessage(terminal);
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }

                // TODO: send the message to service

                // print finished message line
                System.out.printf("\033[2K\rYou wrote: %s\n\r", message);
                // TODO: empty the message
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

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

            // TODO: stop reading and return

            // TODO: delete single character

            controlCharCounter = 0;

            if (character == 27) {
                controlCharCounter = 2;
            }

            // do not display/send control characters
            if (Character.isISOControl(character)) {
                continue;
            }
            // TODO: print character to console and append it to message
        }
    }

    @Override
    public int run(String... args) {
        // TODO: claim name from server

        // TODO: append token in header

        // TODO: attach headers to stub

        // TODO: create output stream

        // TODO: connect to service and start printing incoming messages

        return 0;
    }
}
