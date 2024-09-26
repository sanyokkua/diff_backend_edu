package ua.kostenko.tasks.app.utility;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import ua.kostenko.tasks.app.dto.ResponseDto;

import static org.junit.jupiter.api.Assertions.*;

class ResponseDtoUtilsTest {

    private ResponseDto<String> responseDto;

    @BeforeEach
    void setUp() {
        responseDto = ResponseDto.<String>builder()
                                 .statusCode(200)
                                 .statusMessage("OK")
                                 .data("Sample Data")
                                 .error(null)
                                 .build();
    }

    @Test
    void testResponseDtoBuilder() {
        assertNotNull(responseDto);
        assertEquals(200, responseDto.getStatusCode());
        assertEquals("OK", responseDto.getStatusMessage());
        assertEquals("Sample Data", responseDto.getData());
        assertNull(responseDto.getError());
    }

    @Test
    void testResponseDtoWithNullData() {
        ResponseDto<String> dtoWithNullData =
                ResponseDto.<String>builder().statusCode(200).statusMessage("OK").data(null).error(null).build();
        assertNotNull(dtoWithNullData);
        assertNull(dtoWithNullData.getData());
    }

    @Test
    void testBuildDtoResponse() {
        ResponseEntity<ResponseDto<String>> responseEntity = ResponseDtoUtils.buildDtoResponse("Data", HttpStatus.OK);

        assertEquals(HttpStatus.OK, responseEntity.getStatusCode());
        assertEquals("Data", responseEntity.getBody().getData());
        assertEquals(200, responseEntity.getBody().getStatusCode());
        assertEquals("OK", responseEntity.getBody().getStatusMessage());
    }

    @Test
    void testBuildDtoResponseWithNullData() {
        ResponseEntity<ResponseDto<String>> responseEntity = ResponseDtoUtils.buildDtoResponse(null, HttpStatus.OK);

        assertEquals(HttpStatus.OK, responseEntity.getStatusCode());
        assertNull(responseEntity.getBody().getData());
    }

    @Test
    void testBuildDtoResponseWithCookie() {
        ResponseCookie cookie = ResponseCookie.from("name", "value").build();
        ResponseEntity<ResponseDto<String>> responseEntity =
                ResponseDtoUtils.buildDtoResponse("Data", HttpStatus.OK, cookie);

        assertEquals(HttpStatus.OK, responseEntity.getStatusCode());
        assertEquals("Data", responseEntity.getBody().getData());
        assertEquals(200, responseEntity.getBody().getStatusCode());
        assertEquals("OK", responseEntity.getBody().getStatusMessage());
        assertEquals(cookie.toString(), responseEntity.getHeaders().getFirst(HttpHeaders.SET_COOKIE));
    }

    @Test
    void testBuildDtoErrorResponse() {
        Exception ex = new Exception("Test Exception");
        ResponseEntity<ResponseDto<String>> responseEntity =
                ResponseDtoUtils.buildDtoErrorResponse("Data", HttpStatus.INTERNAL_SERVER_ERROR, ex);

        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR, responseEntity.getStatusCode());
        assertEquals("Data", responseEntity.getBody().getData());
        assertEquals(500, responseEntity.getBody().getStatusCode());
        assertEquals("INTERNAL_SERVER_ERROR", responseEntity.getBody().getStatusMessage());
        assertTrue(responseEntity.getBody().getError().contains("Test Exception"));
    }

    @Test
    void testCreateErrorResponseBody() {
        Exception ex = new Exception("Another Test Exception");
        ResponseDto<String> errorResponse =
                ResponseDtoUtils.createErrorResponseBody("Data", HttpStatus.INTERNAL_SERVER_ERROR, ex);

        assertEquals(500, errorResponse.getStatusCode());
        assertEquals("INTERNAL_SERVER_ERROR", errorResponse.getStatusMessage());
        assertTrue(errorResponse.getError().contains("Another Test Exception"));
        assertEquals("Data", errorResponse.getData());
    }
}