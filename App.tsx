import React, { useState, useEffect } from 'react';
import { Calendar, momentLocalizer } from 'react-big-calendar';
import moment from 'moment';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL;
const localizer = momentLocalizer(moment);

interface Booking {
  id: number;
  title: string;
  start: Date;
  end: Date;
}

interface User {
  id: number;
  name: string;
  membershipType: string;
}

interface Props {
  isStaff: boolean;
}

const FitFleet: React.FC<Props> = ({ isStaff }) => {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        const response = await axios.get(`${API_URL}/bookings`);
        setBookings(response.data);
      } catch (error) {
        console.error('Error fetching bookings', error);
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, []);

  useEffect(() => {
    if (isStaff) {
      const fetchUsers = async () => {
        try {
          const response = await axios.get(`${API_URL}/users`);
          setUsers(response.data);
        } catch (error) {
          console.error('Error fetching users', error);
        }
      };

      fetchUsers();
    }
  }, [isStaff]);

  const handleEventSelect = (event: Booking) => {
    console.log(`Selected booking: ${event.title}`);
  };

  const renderDashboard = () => {
    return (
      <div>
        <h2>User Dashboard</h2>
        <button>Manage My Bookings</button>
        <button>Update Membership</button>
      </div>
    );
  };

  const renderAdminArea = () => {
    if (!isStaff) return null;

    return (
      <div>
        <h2>Administrative Area</h2>
        <button>View All Bookings</button>
        <button>Manage Users</button>
      </div>
    );
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>FitFleet</h1>
      {isStaff && renderAdminArea()}
      {!isStaff && renderDashboard()}
      <Calendar
        localizer={localizer}
        events={bookings}
        startAccessor="start"
        endAccessor="end"
        onSelectEvent={handleEventSelect}
        style={{ height: 500 }}
      />
    </div>
  );
};

export default FitFleet;